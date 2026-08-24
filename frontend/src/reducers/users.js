import { userList } from "constants/users";

const initialState = {
  userList: {},
};

const users = (state = initialState, action) => {
  switch (action.type) {
    case userList:
      return { ...state, userList: action.payload };

    default:
      return state;
  }
};
export default users;
