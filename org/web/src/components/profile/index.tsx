'use client';

import React from 'react';
import { UserProfile, useAuth, useUser, SignIn } from '@clerk/clerk-react';

export const UserProfilePage = () => {
    const { isSignedIn } = useAuth();
    const { user } = useUser(); 

    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-100 dark:bg-navy-900">
            {isSignedIn && user ? (
                <div className="w-full max-w-4xl p-6 bg-white rounded-lg shadow-lg dark:bg-navy-900">
                    <UserProfile
                        additionalOAuthScopes={{
                            google: ['calendar', 'drive'],
                            github: ['repo'],
                        }}
                    />
                </div>
            ) : (
                <div className="text-center text-gray-800 dark:text-white">
                    <SignIn />
                </div>
            )}
        </div>
    );
};

export default UserProfilePage;
