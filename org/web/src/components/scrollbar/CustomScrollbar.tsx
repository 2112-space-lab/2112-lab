import React from "react";
import { Scrollbars } from "react-custom-scrollbars-2";

export const CustomScrollbar = ({
  children,
  style,
}: {
  children: React.ReactNode;
  style?: React.CSSProperties;
}) => {
  const renderTrackVertical = ({ style, ...props }: any) => {
    const trackStyle = {
      position: "absolute",
      width: 6,
      right: 2,
      top: 2,
      bottom: 2,
      borderRadius: 3,
      background: "transparent",
      transition: "opacity 200ms ease 0s",
      opacity: 0,
    };
    return <div style={{ ...style, ...trackStyle }} {...props} />;
  };

  const renderThumbVertical = ({ style, ...props }: any) => {
    const thumbStyle = {
      borderRadius: 3,
      background: "rgba(222, 222, 222, 0.3)",
    };
    return <div style={{ ...style, ...thumbStyle }} {...props} />;
  };

  const renderTrackHorizontal = ({ style, ...props }: any) => {
    const trackStyle = {
      position: "absolute",
      height: 6,
      left: 2,
      right: 2,
      bottom: 2,
      borderRadius: 3,
      background: "transparent",
      transition: "opacity 200ms ease 0s",
      opacity: 0,
    };
    return <div style={{ ...style, ...trackStyle }} {...props} />;
  };

  const renderThumbHorizontal = ({ style, ...props }: any) => {
    const thumbStyle = {
      borderRadius: 3,
      background: "rgba(222, 222, 222, 0.3)",
    };
    return <div style={{ ...style, ...thumbStyle }} {...props} />;
  };

  const renderView = ({ style, ...props }: any) => {
    const viewStyle = {
      marginBottom: -22, // Account for horizontal scrollbar
      marginRight: -22, // Account for vertical scrollbar
    };
    return <div style={{ ...style, ...viewStyle }} {...props} />;
  };

  return (
    <Scrollbars
      style={style}
      renderTrackVertical={renderTrackVertical}
      renderThumbVertical={renderThumbVertical}
      renderTrackHorizontal={renderTrackHorizontal}
      renderThumbHorizontal={renderThumbHorizontal}
      renderView={renderView}
    >
      {children}
    </Scrollbars>
  );
};
