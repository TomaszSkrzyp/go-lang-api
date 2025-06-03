
import './App.css';
import './main.css';
import React, {useEffect,useState} from 'react'


function App() {
  const [tasks, setTasks] = useState([]);
  const [page, setPage] = useState(1);
  const [limit] = useState(5); 
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  const [newTask, setNewTask] = useState('');
  const [newStatus, setNewStatus] = useState('');

  const [editingTaskId, setEditingTaskId] = useState(null);
  const [editText, setEditText] = useState("");
  const [editStatus, setEditStatus] = useState("");

  const statusOptions = ["Pending", "In Progress", "Completed","Canceled"];

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

  
  const moveStatusUp = async (task) => {
    try {
      const response = await fetch(`http://localhost:8090/todos/${task.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ changeUp: "true" }),
      });
      if (!response.ok) throw new Error("Failed to move status up");
      fetchTasks(page);
    } catch (err) {
      console.error(err);
    }
  };

  const removeTask = async (taskId) => {
    try {
      const response = await fetch(`http://localhost:8090/todos/${taskId}`, {
        method: 'DELETE',
      });
      if (!response.ok) throw new Error("Failed to remove task");

      if (tasks.length === 1 && page > 1) {
        setPage(page - 1);
      } else {
        fetchTasks(page);
      }
    } catch (err) {
      console.error(err);
    }
  };

 const updateTask = async (id) => {
  try {
    const res = await fetch(`http://localhost:8090/todos/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        task: editText,
        status: editStatus,
      }),
    });

    if (!res.ok) throw new Error("Failed to update task");

    await fetchTasks(page);
    setEditingTaskId(null);
  } catch (err) {
    console.error(err);
  }
};

  
  const addTask = async (e) => {
    e.preventDefault();
    if (!newTask.trim()) {
      alert("Task name is required.");
      return;
    }

    try {
      await fetch('http://localhost:8090/todos', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ task: newTask, status: newStatus || undefined }),
      });

      setNewTask('');
      setNewStatus('');
      fetchTasks(page);
    } catch (err) {
      console.error(err);
    }
  };
  const totalPages = Math.ceil(total / limit);

return (<div style={{ padding: '20px' }}>
  <h1>TODO List (Page {page})</h1>

  <form onSubmit={addTask} style={{ marginBottom: '20px' }}>
    <input
      type="text"
      placeholder="Task description"
      value={newTask}
      onChange={(e) => setNewTask(e.target.value)}
      required
    />

    <select
      value={newStatus}
      onChange={(e) => setNewStatus(e.target.value)}
    >
      <option value="">Select Status (optional)</option>
      {statusOptions.map((status) => (
        <option key={status} value={status}>
          {status}
        </option>
      ))}
    </select>

    <button type="submit">Add Task</button>
  </form>

  {loading && <p>Loading...</p>}

  {!loading && (
    <ul>
      {tasks.map((task) => (
        <li
          key={task.id}
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          
          {editingTaskId === task.id ? (
            <>
              <div>
                <input
                  type="text"
                  value={editText}
                  onChange={(e) => setEditText(e.target.value)}
                />
                <select
                  value={editStatus}
                  onChange={(e) => setEditStatus(e.target.value)}
                >
                  {statusOptions.map((s) => (
                    <option key={s} value={s}>
                      {s}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <button
                  onClick={() => updateTask(task.id)}>
                  Save
                </button>
                <button onClick={() => setEditingTaskId(null)}>Cancel</button>
              </div>
            </>
          ) : (
            <>
              <div>
                {task.task} - <em>{task.status}</em>
              </div>
              <div>
                <button
                  onClick={() => moveStatusUp(task)}
                  disabled={task.status === 'Completed'}
                >
                  Status Up
                </button>
                <button
                  onClick={() => {
                    setEditingTaskId(task.id);
                    setEditText(task.task);
                    setEditStatus(task.status);
                  }}
                >
                  Edit
                </button>
                <button
                  onClick={() => removeTask(task.id)}
                  style={{ backgroundColor: 'red', color: 'white' }}
                >
                  Remove
                </button>
              </div>
            </>
          )}
        </li>
      ))}
    </ul>
  )}

  <div style={{ marginTop: 20 }}>
    <button onClick={() => setPage((p) => Math.max(p - 1, 1))} disabled={page === 1}>
      Prev
    </button>

    <span style={{ margin: '0 10px' }}>
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
