import React, { useEffect, useMemo } from "react";
import { useQuery, useSubscription, useMutation, gql } from "@apollo/client";
import { Spinner, Box, Text, Center, IconButton } from "@chakra-ui/react";
import VisibilityTimeline from "./VisibilityTimeline";
import { FiRefreshCw } from "react-icons/fi";

const GET_CACHED_VISIBILITIES = gql`
  query GetCachedSatelliteVisibilities(
    $uid: String!
    $userLocation: UserLocationInput!
    $startTime: String!
    $endTime: String!
  ) {
    cachedSatelliteVisibilities(
      uid: $uid
      userLocation: $userLocation
      startTime: $startTime
      endTime: $endTime
    ) {
      satelliteId
      satelliteName
      aos
      los
    }
  }
`;

const VISIBILITY_SUBSCRIPTION = gql`
  subscription OnVisibilityUpdated(
    $uid: String!
    $userLocation: UserLocationInput!
    $startTime: String!
    $endTime: String!
  ) {
    satelliteVisibilityUpdated(
      uid: $uid
      userLocation: $userLocation
      startTime: $startTime
      endTime: $endTime
    ) {
      satelliteId
      satelliteName
      aos
      los
    }
  }
`;

const REQUEST_VISIBILITIES_MUTATION = gql`
  mutation RequestSatelliteVisibilities(
    $uid: String!
    $userLocation: UserLocationInput!
    $startTime: String!
    $endTime: String!
  ) {
    requestSatelliteVisibilities(
      uid: $uid
      userLocation: $userLocation
      startTime: $startTime
      endTime: $endTime
    )
  }
`;

interface VisibilityTimelineWithDataProps {
  uid: string;
  userLocation: { latitude: number; longitude: number };
}

const VisibilityTimelineWithData: React.FC<VisibilityTimelineWithDataProps> = ({ uid, userLocation }) => {
  const radius = 1; // 500 km radius
  const horizon = 30; // 30 degrees horizon angle

  // Set up time range
  const startTime = useMemo(() => new Date().toISOString(), []);
  const endTime = useMemo(() => new Date(Date.now() + 60 * 60 * 1000 * 24).toISOString(), []);

  const userLocationInput = {
    ...userLocation,
    radius,
    horizon,
    uid,
  };

  // Mutation to request satellite visibilities
  const [requestVisibilities, { loading: mutationLoading, error: mutationError }] = useMutation(
    REQUEST_VISIBILITIES_MUTATION
  );

  // Query to fetch initial cached data
  const { data, loading: queryLoading, error: queryError } = useQuery(GET_CACHED_VISIBILITIES, {
    variables: { uid, userLocation: userLocationInput, startTime, endTime },
    skip: !uid,
  });

  // Subscription for real-time updates
  const { data: subscriptionData, error: subscriptionError } = useSubscription(VISIBILITY_SUBSCRIPTION, {
    variables: { uid, userLocation: userLocationInput, startTime, endTime },
    skip: !uid,
  });
  useEffect(() => {
    if (uid) {
      requestVisibilities({
        variables: { uid, userLocation: userLocationInput, startTime, endTime },
      }).catch((err) => {
        console.error("Error requesting satellite visibilities:", err);
      });
    }
  }, [uid]);

  // Combine query and subscription data
  const visibilities = useMemo(() => {
    const queriedData = data?.cachedSatelliteVisibilities || [];
    const liveUpdates = subscriptionData?.satelliteVisibilityUpdated || [];
    return [...queriedData, ...liveUpdates];
  }, [data, subscriptionData]);

  if (queryLoading || mutationLoading) {
    return (
      <Center>
        <Spinner />
        <Text ml={2}>Loading visibilities...</Text>
      </Center>
    );
  }
  // Handle error states
  if (queryError || subscriptionError || mutationError) {
    console.error("Query Error:", queryError);
    console.error("Subscription Error:", subscriptionError);
    console.error("Mutation Error:", mutationError);
    return (
      <Box>
        <Text color="red.500">
          Error loading visibilities:
          {queryError?.message || subscriptionError?.message || mutationError?.message}
        </Text>
      </Box>
    );
  }

  // Combine and deduplicate visibility data
  const timelineData = visibilities
    .map((visibility) => {
      const aosDate = new Date(visibility.aos);
      const losDate = new Date(visibility.los);
      return {
        day: aosDate.toUTCString().split(",")[1]?.trim().split(" ")[0] || "N/A",
        weekday: aosDate.toUTCString().split(",")[0] || "N/A",
        hours: `${aosDate.toUTCString().split(" ")[4]} - ${losDate.toUTCString().split(" ")[4]}`,
        title: "Visibility Window",
        satellite: visibility.satelliteName || "Unknown Satellite",
        spaceID: visibility.satelliteId || "N/A",
        aosDate,
      };
    }).sort((a, b) => a.aosDate.getTime() - b.aosDate.getTime());
  return <VisibilityTimeline location={userLocation} data={timelineData} />;
};

export default VisibilityTimelineWithData;
