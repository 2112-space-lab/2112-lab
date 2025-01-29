'use client';
import { useUser } from '@clerk/nextjs';
import { UserProfilePage } from 'components/profile/index';


const ProfileSettingsPage = () => {
    const { isLoaded, isSignedIn, user } = useUser();

    if (!isLoaded) {
        return <div>Loading...</div>;
    }

    return (
        <UserProfilePage></UserProfilePage>
    );
};

export default ProfileSettingsPage;
