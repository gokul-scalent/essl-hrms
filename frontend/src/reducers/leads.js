import { leadList } from "constants/leads";

const initialState = {
  leadList: {},
};

const leads = (state = initialState, action) => {
  switch (action.type) {
    case leadList:
      return { ...state, leadList: action.payload };

    default:
      return state;
  }
};
export default leads;
