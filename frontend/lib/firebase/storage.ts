import { initializeApp, getApps } from "firebase/app"
import { getStorage } from "firebase/storage"
import firebaseConfig from "./config";

// Initialize a new app if it doesn't exist
if (getApps().length === 0) {
    try {
        // Initialize Firebase
        initializeApp(firebaseConfig);
    } catch (error) {
        console.error("Firebase initialization error: ", error.stack);
    }
}

export default getStorage();