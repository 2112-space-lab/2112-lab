'use client';

import { useEffect, useState } from 'react';
import SearchTableUsers from 'components/admin/main/users/users-overview/SearchTableUsersOverivew';
import useUserServiceStore from 'services/userService';
import Card from 'components/card';

const UserOverview = () => {
  // const { users, fetchPaginatedUsers, loading, error } = useUserServiceStore();
  // const [tableData, setTableData] = useState([]);

  // useEffect(() => {
  //   const fetchData = async () => {
  //     await fetchPaginatedUsers(0, 10);
  //   };
  //   fetchData();
  // }, [fetchPaginatedUsers]);

  // useEffect(() => {
  //   const transformedData = users.map((user) => ({
  //     name: user.name,
  //     email: user.email,
  //     username: user.username,
  //     date: user.date,
  //     type: user.type,
  //     actions: 'Actions',
  //   }));
  //   setTableData(transformedData);
  // }, [users]);

  return (
    <div></div>
    // <Card extra={'w-full h-full mt-3'}>
    //   {loading ? (
    //     <div>Loading...</div>
    //   ) : error ? (
    //     <div>Error: {error}</div>
    //   ) : (
    //     <SearchTableUsers tableData={tableData} />
    //   )}
    // </Card>
  );
};

export default UserOverview;
