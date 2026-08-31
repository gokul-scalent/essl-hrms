import { combineReducers } from "redux";
import users from "./users";
import employee from "./employee";
import attendanceLog from "./attendanceLog";

export const reducers = combineReducers({
  users,
  employee,
  attendanceLog,
});
