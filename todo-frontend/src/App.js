
import './App.css';
import React, {useEffect,useState} from 'react'

function App() {
  const [tasks, setTasks] = useState([]);
  const [page, setPage] = useState(1);
  const [limit] = useState(5); // tasks per page
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  const fetchTasks = async (page) => {
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8090/todos?page=${page}&limit=${limit}`);
      if (!response.ok) throw new Error("Failed to fetch");

      const data = await response.json();
      setTasks(data.tasks);
      setTotal(data.total);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTasks(page);
  }, [page]);

  const totalPages = Math.ceil(total / limit);

  return (
    <div>
      <h1>TODO List (Page {page})</h1>
      {loading && <p>Loading...</p>}
      {!loading && (
        <ul>
          {tasks.map((task) => (
            <li key={task.id}>
              {task.task} - <em>{task.status}</em>
            </li>
          ))}
        </ul>
      )}

      <div style={{ marginTop: 20 }}>
        <button onClick={() => setPage((p) => Math.max(p - 1, 1))} disabled={page === 1}>
          Prev
        </button>

        <span style={{ margin: "0 10px" }}>
          Page {page} of {totalPages}
        </span>

        <button onClick={() => setPage((p) => Math.min(p + 1, totalPages))} disabled={page === totalPages}>
          Next
        </button>
      </div>
    </div>
  );
}

export default App;
