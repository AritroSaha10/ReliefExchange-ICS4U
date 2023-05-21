import { DocumentReference } from "firebase/firestore"

/**
 * Data schema for User Data, separate from the user data directly from our authentication server.
 */
export default interface UserData {
    display_name: string,
    email: string,
    registered_date: string,
    admin: string,
    posts: DocumentReference[],
    donations_made: Number
}