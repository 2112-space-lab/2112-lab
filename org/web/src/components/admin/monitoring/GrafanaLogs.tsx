import React from "react";
import { Box } from "@chakra-ui/react";

export const MonitoringGrafanaLogsPage: React.FC = () => {

    return (
        <Box p={4} w="100%" h="100%">
            <iframe
                src="http://localhost:3001/grafana/d/f494b60b-5027-43c5-8f59-9a2874b8a9ee/logs?orgId=1&from=1737740520442&to=1737740820442"
                width="100%"
                height="1000px"
            ></iframe>
        </Box>
    );
};

export default MonitoringGrafanaLogsPage;
