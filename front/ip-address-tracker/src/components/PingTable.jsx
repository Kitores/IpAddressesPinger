import React, { useEffect, useState } from 'react';
import axios from 'axios';

const PingTable = () => {
    const [data, setData] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('http://backend:8080/getListIp');
                console.log('Response status:', response.status);
                if (!response.ok) {
                    throw new Error('Network response was not ok: ${response.statusText}');
                }
                const jsonData = await response.json();
                setData(jsonData);
            } catch (error) {
                console.error('Ошибка при получении данных:', error);
            }
        };

        fetchData();
        const intervalId = setInterval(fetchData, 30000);

        return () => clearInterval(intervalId);
    }, []);

    return (
        <table>
            <thead>
            <tr>
                <th>IP Address</th>
                <th>Ping Time</th>
                {/*<th>Last Success Date</th>*/}
                <th>Status</th>
            </tr>
            </thead>
            <tbody>
            {Array.isArray(data) && data.map((ping, index) => (
                <tr key={index}>
                    <td>{ping.ip}</td>
                    <td>{new Date(ping.pingTime).toLocaleString()}</td>
                    <td>{ping.success}</td>
                    {/*<td>{new Date(ping.lastSuccessDate).toLocaleString()}</td>*/}
                </tr>
            ))}
            </tbody>
        </table>
    );
};

export default PingTable;
