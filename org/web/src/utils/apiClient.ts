import axios from "axios";
import { useSession } from "@clerk/nextjs";

// Helper function to make authenticated API calls using Clerk's `getToken`
// const apiCallWithAuth = async () => {
//     const session = useSession(); // Access the session object
//     const token = await session.session?.getToken(); // Fetch the current user's token dynamically
//     if (!token) {
//         console.error("Failed to retrieve token. User might not be signed in.");
//         return;
//     }

//     const response = await fetch("/api/endpoint", {
//         method: "GET",
//         headers: {
//             Authorization: `Bearer ${token}`,
//         },
//     });

//     const data = await response.json();
//     console.log(data);
// };

// Create Axios instance
const apiClient = axios.create({
    baseURL: process.env.API_BASE_URL || "http://localhost:8081", // Fallback for local development
    headers: { Accept: "application/json", "Content-Type": "application/json" },
    timeout: 30000, // Timeout in milliseconds
});

// Request interceptor for dynamic token injection
apiClient.interceptors.request.use(
    async (config) => {
        const session = useSession(); // Access Clerk session
        if (session.isLoaded && session.isSignedIn) {
            const token = await session.session?.getToken(); // Fetch token dynamically
            if (token) {
                config.headers.Authorization = `Bearer ${token}`;
            }
        } else {
            console.warn("Session not loaded or user not signed in. Skipping Authorization header.");
        }
        return config;
    },
    (error) => {
        console.error("Request setup error:", error);
        return Promise.reject(error);
    }
);

// Response interceptor
apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response) {
            switch (error.response.status) {
                case 401:
                    console.warn("Unauthorized! Redirecting to login...");
                    break;
                case 406:
                    console.error("HTTP 406: Not Acceptable. Check headers or payload.");
                    break;
                case 500:
                    console.error(
                        "HTTP 500: Internal Server Error:",
                        error.response.data.message || "No details provided."
                    );
                    break;
                default:
                    console.warn(
                        `HTTP ${error.response.status}: ${error.response.data.message || "No details provided."}`
                    );
            }
        } else if (error.request) {
            console.error("Network error or no response received.");
        } else {
            console.error("Request setup error:", error.message);
        }
        return Promise.reject(error);
    }
);

// Debugging logs in development mode
if (process.env.NODE_ENV === "development") {
    apiClient.interceptors.request.use((config) => {
        console.log("Request Config:", config);
        return config;
    });
    apiClient.interceptors.response.use(
        (response) => {
            console.log("Response Data:", response.data);
            return response;
        },
        (error) => {
            console.error("Error Response:", error.response || error);
            return Promise.reject(error);
        }
    );
}

export default apiClient;
