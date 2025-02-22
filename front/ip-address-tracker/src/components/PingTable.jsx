import React, { useEffect, useState } from 'react';
import './styles.css';

const PingTable = () => {
    const [data, setData] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('/getListIp');
                console.log('Response status:', response.status);
                if (!response.ok) {
                    throw new Error(`Network response was not ok: ${response.statusText}`);
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
        <div className="table-container">
            <table>
                <thead>
                <tr>
                    <th>IP Address</th>
                    <th>Ping Time</th>
                    <th>Status</th>
                </tr>
                </thead>
                <tbody>
                {Array.isArray(data) && data.map((ping, index) => (
                    <tr key={index}>
                        <td>{ping.IpAddress}</td>
                        <td>{new Date(ping.PingTime).toLocaleString()}</td>
                        <td>{ping.Status}</td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default PingTable;
