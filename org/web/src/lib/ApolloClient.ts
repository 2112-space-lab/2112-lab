// lib/apolloClient.ts

import { ApolloClient, InMemoryCache } from '@apollo/client';

const AppoloClient = new ApolloClient({
  uri: 'http://localhost:4000/query',
  cache: new InMemoryCache()
});

export default AppoloClient;