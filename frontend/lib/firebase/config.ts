// The configuration data for the firebase project.
const firebaseConfig = {
    apiKey: "AIzaSyDpTBi9KDh0jJvnQqxRiNJGZBVJQy8VZP4",
    authDomain: process.env.NODE_ENV === "production" ? "reliefexchange.aritrosaha.ca" : "ics4u-project.firebaseapp.com",
    projectId: "ics4u-project",
    storageBucket: "ics4u-project.appspot.com",
    messagingSenderId: "42306297294",
    appId: "1:42306297294:web:38997c053c7d1047761f9f",
    measurementId: "G-B1WYFG9DDL"
};

export default firebaseConfig;