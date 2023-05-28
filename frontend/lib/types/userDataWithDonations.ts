import UserData from "./userData"
import Donation from "./donation"

/**
 * Data schema for User Data with Donations, separate from the user data directly from our authentication server.
 * This is different from UserData by having actual Donation data instead of references to them.
 */
export default interface UserDataWithDonations {
    display_name: string,
    email: string,
    registered_date: string,
    admin: string,
    posts: Donation[],
    donations_made: Number
}