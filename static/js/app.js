document.getElementById('registerForm').addEventListener('submit', async (event) => {
    event.preventDefault();
    const username = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;

    const response = await fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
    });

    if (response.ok) {
        alert('Registration successful!');
        window.location.href = '/dashboard';  // Redirect to dashboard
    } else {
        alert('Registration failed!');
    }
});

document.getElementById('loginForm').addEventListener('submit', async (event) => {
    event.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    const response = await fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
    });

    if (response.ok) {
        const data = await response.json(); // Assuming the response returns a JSON object
        const token = data.token; // Get the token from the response

        // Set the token as a cookie
        document.cookie = `jwtToken=${token}; path=/; Secure; SameSite=Strict`;

        alert('Login successful!');
        window.location.href = '/dashboard';  // Redirect to dashboard
    } else {
        alert('Login failed!');
    }
});

// Example of how to access the dashboard
async function accessDashboard() {
    const token = getCookie('jwtToken'); // Retrieve the token from the cookie
    
    const response = await fetch('/dashboard', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}` // Include the token in the Authorization header
        }
    });

    if (response.ok) {
        const dashboardData = await response.text(); // Assuming a simple text response
        console.log(dashboardData); // Handle the dashboard data
    } else {
        alert('Failed to access dashboard!');
    }
}

// Helper function to get a cookie by name
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

accessDashboard();