function connect() {
    const goload = new WebSocket("ws://localhost:{{.Port}}/ws");

    goload.onopen = () => {
        console.log("Connected to GoLoad server");
    };

    goload.onmessage = () => location.reload(true);

    goload.onclose = () => {
        console.log("GoLoad connection closed, retrying...");
        setTimeout(connect, 2000);
    };

    goload.onerror = (error) => console.error("GoLoad error:", error);
}
connect();
