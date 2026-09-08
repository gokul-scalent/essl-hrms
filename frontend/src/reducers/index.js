import { combineReducers } from "redux";
import users from "./users";
import employee from "./employee";
import attendanceLog from "./attendanceLog";
import role from "./role";

export const reducers = combineReducers({
  users,
  employee,
  attendanceLog,
  role,
});
