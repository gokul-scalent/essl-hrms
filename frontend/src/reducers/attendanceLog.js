import {
  attendanceLogList,
  dailyAttendanceLogList,
} from "constants/attendanceLog";

const initialState = {
  attendanceLogList: {},
  dailyAttendanceLogList: {},
};

const attendanceLog = (state = initialState, action) => {
  switch (action.type) {
    case attendanceLogList:
      return { ...state, attendanceLogList: action.payload };

    case dailyAttendanceLogList:
      return { ...state, dailyAttendanceLogList: action.payload };

    default:
      return state;
  }
};
export default attendanceLog;
