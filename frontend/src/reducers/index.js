import { combineReducers } from "redux";
import email from "./email";
import leads from "./leads";
import users from "./users";

export const reducers = combineReducers({
  email,
  leads,
  users,
});
