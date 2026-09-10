import { roleList } from "constants/role";


const initialState = {
  roleList: {},
};

const role = (state = initialState, action) => {
  switch (action.type) {
    case roleList:
      return { ...state, roleList: action.payload };
    default:
      return state;
  }
};

export default role;