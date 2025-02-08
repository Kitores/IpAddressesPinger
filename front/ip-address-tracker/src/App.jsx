import React from 'react';
import IpTable from './components/IpTable.jsx';
import { useState } from 'react';
import './styles.css';

// function App() {
//     return (
//         <div className="App">
//             <IpTable />
//         </div>
//     );
// }

function App() {
    const [data, setData] = useState(null);

    const fetchData = async () => {
        try {
            const response = await fetch('http://localhost:8080/ping');
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            const result = await response.json();
            setData(result);
        } catch (error) {
            console.error('Ошибка:', error);
        }
    };

    return (
        <div className="container">
            <button className="styled-button" onClick={fetchData}>
                Получить данные
            </button>
            {data && (
                <pre>
                    {JSON.stringify(data, null, 2)}
                </pre>
            )}
        </div>
    );
}
export default App;