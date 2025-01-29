/** @type {import('next').NextConfig} */
const path = require('path');
const process = require('process');
const CopyWebpackPlugin = require('copy-webpack-plugin');

// Bundle analyzer for debugging build size
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
});

const pathBuilder = (subpath) => path.join(process.cwd(), subpath);

const nextConfig = withBundleAnalyzer({
  poweredByHeader: false, // Remove the "X-Powered-By: Next.js" header
  trailingSlash: true, // Add trailing slash to all routes
  reactStrictMode: true, // Enable React strict mode
  swcMinify: true, // Use SWC for minification
  basePath: process.env.NEXT_PUBLIC_BASE_PATH || '', // Set the base path if provided
  assetPrefix: process.env.NEXT_PUBLIC_BASE_PATH || '', // Prefix for assets if needed
  images: {
    domains: [
      'images.unsplash.com',
      'i.ibb.co',
      'scontent.fotp8-1.fna.fbcdn.net',
      'picsum.photos',
      'shipixen.com',
    ],
    unoptimized: false, // Optimize images when possible
  },
  env: {
    MAPBOX_TOKEN: process.env.MAPBOX_TOKEN || '', // Use an environment variable for the Mapbox token
  },
  webpack: (config, { webpack, isServer }) => {
    if (!isServer) {
      // Add Cesium-specific configurations for the client-side
      config.plugins.push(
        new CopyWebpackPlugin({
          patterns: [
            {
              from: pathBuilder('node_modules/cesium/Build/Cesium/Workers'),
              to: 'public/cesium/Workers',
              info: { minimized: true },
            },
            {
              from: pathBuilder('node_modules/cesium/Build/Cesium/ThirdParty'),
              to: 'public/cesium/ThirdParty',
              info: { minimized: true },
            },
            {
              from: pathBuilder('node_modules/cesium/Build/Cesium/Assets'),
              to: 'public/cesium/Assets',
              info: { minimized: true },
            },
            {
              from: pathBuilder('node_modules/cesium/Build/Cesium/Widgets'),
              to: 'public/cesium/Widgets',
              info: { minimized: true },
            },
          ],
        }),
        new webpack.DefinePlugin({
          CESIUM_BASE_URL: JSON.stringify('/cesium'), // Define Cesium base URL
        })
      );
    }

    // Add rule for .svg files
    config.module.rules.push({
      test: /\.svg$/,
      use: ['@svgr/webpack'],
    });

    config.module.rules.push({
      test: /\.(mp4|webm)$/,
      type: 'asset/resource',
      generator: {
        filename: 'static/media/[name].[hash][ext]',
      },
    });

    return config;
  },
  output: 'standalone',
});

module.exports = nextConfig;
