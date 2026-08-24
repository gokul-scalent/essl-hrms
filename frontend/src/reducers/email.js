import { emailList } from "constants/email";

const initialState = {
  emailList: {},
};

const email = (state = initialState, action) => {
  switch (action.type) {
    case emailList:
      return { ...state, emailList: action.payload };

    default:
      return state;
  }
};
export default email;
