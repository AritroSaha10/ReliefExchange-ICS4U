import Donation from "./donation"
import UserData from "./userData"

/**
 * Data schema for a Donation combined with it's owner's user data. Not available
 * directly from the database, but both need to be manually combined.
 */
export default interface DonationWithUserData extends Donation {
    owner: UserData
}