package schedule_manager

import (
	"bytes"
	"context"
	"encoding/json"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/NEO-TAT/tat_auto_roll_call_service/src/constants"
	"github.com/NEO-TAT/tat_auto_roll_call_service/src/types"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
)

// Contains checks if the collection contains an element equal to [element].
//
// This operation will check each element in order for being equal to
// `element`, unless it has a more efficient way to find an element
// equal to `element`.
//
// The equality used to determine whether `element` is equal to an element of
// the iterable defaults to the `==` of the element.
func Contains[E comparable](collection []E, element E) bool {
	for _, v := range collection {
		if v == element {
			return true
		}
	}

	return false
}

type GetRollCallResult struct {
	Status   bool   `json:"status,omitempty"`
	Msg      string `json:"msg,omitempty"`
	RollCall *struct {
		Id       string `json:"id,omitempty"`
		CourseId string `json:"course_id,omitempty"`
		Name     string `json:"name,omitempty"`
		IsGPS    string `json:"is_gps,omitempty"`
		Record   struct {
			Answered bool `json:"answered,omitempty"`
		}
	} `json:"rollcall,omitempty"`
}

var notifiedRollCallIds []string

// ParseJSONBody parses http request or response application/json data to a declared target instance.
func ParseJSONBody(source io.ReadCloser, target any) error {
	// un-serialize the origin source.
	dec := json.NewDecoder(source)

	// prevent unknown field appear in decode result.
	//dec.DisallowUnknownFields()

	// put the decoded result to target(reference).
	err := dec.Decode(&target)
	if err != nil {
		return err
	}

	// verify the data is clear.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return err
	}

	return nil
}

func (m manager) UpdateSchedules() (err error) {
	ctx := context.Background()
	userCollectionRef := m.fb.Store.Collection(constants.KAutoRollCallUsersCollectionName)

	// Get All users.
	userSnapShots, err := userCollectionRef.Documents(ctx).GetAll()
	if err != nil {
		m.logger.Error(err)
		return err
	}

	targetCourseWithScheduleIds := make(map[string][]string)
	scheduleIdWithDeviceTokens := make(map[string]string)
	schedules := make(map[string]map[string]interface{})

	for _, userSnapShot := range userSnapShots {
		scheduleUserInfo := &types.AutoRollCallUser{}
		err := mapstructure.Decode(userSnapShot.Data(), scheduleUserInfo)
		if err != nil {
			m.logger.Error(err)
			return err
		}

		scheduleCollectionRef := userSnapShot.Ref.Collection(constants.KSchedulesFieldName)
		scheduleSnapShots, err := scheduleCollectionRef.Documents(ctx).GetAll()
		if err != nil {
			m.logger.Error(err)
			return err
		}

		for _, scheduleSnapShot := range scheduleSnapShots {
			schedule := scheduleSnapShot.Data()
			if !(schedule["enabled"].(bool)) {
				continue
			}

			timeRange := schedule["time_range"].(map[string]interface{})
			weekDay := timeRange["selected_week_day"].(string)
			if weekDay != "thursday" {
				continue
			}
			//period := timeRange["period"].(map[string]interface{})
			//
			//startTime := period["start_time"].(map[string]interface{})
			//endTime := period["end_time"].(map[string]interface{})
			//
			//startTimeHour := startTime["hour"].(string)
			//endTimeHour := endTime["hour"].(string)
			//startTimeMinute := startTime["minute"].(string)
			//endTimeMinute := endTime["minute"].(string)
			//
			//currentTime := time.Now()
			//currentHour := currentTime.Hour()
			//currentMinute := currentTime.Minute()

			//err := mapstructure.Decode(scheduleSnapShot.Data(), schedule)
			//if err != nil {
			//	m.logger.Error(err)
			//	return err
			//}

			schedules[schedule["id"].(string)] = schedule

			targetCourse := schedule["target_course"].(map[string]interface{})
			targetCourseId := targetCourse["id"].(string)

			targetCourseSchedules := targetCourseWithScheduleIds[targetCourseId]
			targetCourseWithScheduleIds[targetCourseId] = append(targetCourseSchedules, schedule["id"].(string))

			scheduleIdWithDeviceTokens[schedule["id"].(string)] = scheduleUserInfo.UserDeviceTokens[0]
		}
	}
	//
	//m.logger.Infoln(targetCourseWithScheduleIds)
	//m.logger.Infoln(scheduleIdWithDeviceTokens)

	m.cronCheckHasRollCall(schedules, targetCourseWithScheduleIds, scheduleIdWithDeviceTokens)

	return nil
}

type msgCredential struct {
	deviceToken string
	courseName  string
}

func (m manager) cronCheckHasRollCall(schedules map[string]map[string]interface{}, targetCourseWithScheduleIds map[string][]string, scheduleIdWithDeviceTokens map[string]string) {
	var shouldNotifyList []msgCredential

	for courseId, scheduleIds := range targetCourseWithScheduleIds {
		for _, scheduleId := range scheduleIds {
			targetUser := schedules[scheduleId]["target_user"].(map[string]interface{})
			accessToken := targetUser["access_token"].(string)
			userId := targetUser["id"].(string)

			rollCallExists := m.checkIfRollCallExists(courseId, accessToken, userId)
			if rollCallExists {
				deviceToken := scheduleIdWithDeviceTokens[scheduleId]
				courseName := schedules[scheduleId]["target_course"].(map[string]interface{})["name"].(string)

				shouldNotifyList = append(shouldNotifyList, msgCredential{
					deviceToken: deviceToken,
					courseName:  courseName,
				})
			}
		}
	}

	m.logger.Warningln(shouldNotifyList)

	for _, shouldNotify := range shouldNotifyList {
		m.sendNotification(shouldNotify.deviceToken, shouldNotify.courseName)
	}
}

func (m manager) checkIfRollCallExists(courseID, accessToken, userId string) (ok bool) {
	const getRollCallUrl = "https://irs.zuvio.com.tw/app_v2/getRollcall"
	jsonBody := []byte(fmt.Sprintf(`{"user_id": "%s", "accessToken": "%s", "course_id": "%s"}`, userId, accessToken, courseID))
	bodyReader := bytes.NewReader(jsonBody)
	res, err := http.Post(getRollCallUrl, "application/json", bodyReader)
	if err != nil {
		m.logger.Errorln(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			m.logger.Errorln(err)
		}
	}(res.Body)

	result := &GetRollCallResult{}
	err = ParseJSONBody(res.Body, result)
	if err != nil {
		fmt.Println(err)
	}

	ok = result.RollCall != nil && !result.RollCall.Record.Answered && !Contains(notifiedRollCallIds, result.RollCall.Id)
	//if ok {
	//	notifiedRollCallIds = append(notifiedRollCallIds, result.RollCall.Id)
	//}

	return
}

func (m manager) sendNotification(deviceToken string, courseName string) {
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Token: deviceToken,
		Notification: &messaging.Notification{
			Title: "Roll Call Remind 點名提醒～",
			Body:  fmt.Sprintf("快打開 TAT，%s 開放點名中～", courseName),
		},
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := m.fb.Message.Send(context.Background(), message)
	if err != nil {
		m.logger.Errorln(err)
	}
	// Response is a message ID string.
	m.logger.Infoln("Successfully sent message:", response)
}
