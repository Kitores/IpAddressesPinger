import React, { useState, useEffect } from 'react';
import axios from 'axios';
import Table from 'react-bootstrap/Table';

function IpTable() {
    const [data, setData] = useState([]);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await axios.get('http://your-api-endpoint.com/ips');
                setData(response.data);
            } catch (error) {
                console.error('Ошибка при получении данных:', error);
            }
        }

        fetchData();
    }, []);

    return (
        <div className="container mt-5">
            <h2>Список IP-адресов</h2>
            <Table>
                <thead>
                <tr>
                    <th>IP Address</th>
                    <th>Location</th>
                </tr>
                </thead>
                <tbody>
                {data.map((item, index) => (
                    <tr key={index}>
                        <td>{item.ip}</td>
                        <td>{item.description}</td>
                    </tr>
                ))}
                </tbody>
            </Table>
        </div>
    );
}

export default IpTable;
