import { attendanceLogList } from "constants/attendanceLog";
import * as API from "../API/index";
import { authorizedUser } from "components/CommonComponent/CommonFunction";

//get users list
export const getAttendanceLogList = async (dispatch,pageNum,filters,searchString,) => {
  try {
    const res = await API.get_attendance_log_list(pageNum,filters,searchString,);
    dispatch({ type: attendanceLogList, payload: res.data });
  } catch (error) {
    authorizedUser(error.response?.data);
    dispatch({ type: attendanceLogList, payload: error?.response?.data || error });
  }
};