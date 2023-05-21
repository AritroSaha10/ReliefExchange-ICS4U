import { initializeApp, getApps } from "firebase/app"
import { getAuth } from "firebase/auth"

// Initialize a new app if it doesn't exist
if (getApps().length === 0) {
    try {
        // Initialize Firebase
        initializeApp({
            apiKey: "AIzaSyDpTBi9KDh0jJvnQqxRiNJGZBVJQy8VZP4",
            authDomain: "ics4u-project.firebaseapp.com",
            projectId: "ics4u-project",
            storageBucket: "ics4u-project.appspot.com",
            messagingSenderId: "42306297294",
            appId: "1:42306297294:web:38997c053c7d1047761f9f",
            measurementId: "G-B1WYFG9DDL"
        });
    } catch (error) {
        console.error("Firebase initialization error: ", error.stack);
    }
}

export default getAuth();