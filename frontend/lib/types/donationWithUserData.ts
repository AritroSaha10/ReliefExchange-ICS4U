import Donation from "./donation"
import UserData from "./userData"

export default interface DonationWithUserData extends Donation {
    owner: UserData
}