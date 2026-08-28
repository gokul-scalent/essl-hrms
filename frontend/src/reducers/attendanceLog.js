import { attendanceLogList } from "constants/attendanceLog";


const initialState = {
  attendanceLogList: {},
};

const attendanceLog = (state = initialState, action) => {
  switch (action.type) {
    case attendanceLogList:
      return { ...state, attendanceLogList: action.payload };

    default:
      return state;
  }
};
export default attendanceLog;
