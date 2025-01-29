import React from "react";
import { Box } from "@chakra-ui/react";

export const MonitoringGrafanaHealthServicePage: React.FC = () => {

    return (
        <Box p={4} w="100%" h="100%">
            <iframe
                src="http://localhost:3001/grafana/d/xtkCtBkiz/prometheus-blackbox-exporter?orgId=1&refresh=10s&from=1737736823013&to=1737740423013"
                width="100%"
                height="1000px"
            ></iframe>
        </Box>
    );
};

export default MonitoringGrafanaHealthServicePage;
