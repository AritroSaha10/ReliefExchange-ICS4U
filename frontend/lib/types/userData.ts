import { DocumentReference } from "firebase/firestore"

export default interface UserData {
    display_name: string,
    email: string,
    registered_date: string,
    admin: string,
    posts: DocumentReference[],
    donations_made: Number
}