window.addEventListener("load", () => {
    document.querySelector("button#login").addEventListener("click", () => {
        const username = document.querySelector("#username").value;
        const password = document.querySelector("#password").value;
        const baseUrl = `${window.location.protocol}//${window.location.host}`;
        fetch(`${baseUrl}/api/login`, {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({ username, password })
        })
            .then(res => {
                window.location.assign(`${baseUrl}/index.html`);
            })
            .catch(err => {
                document.querySelector("#result").innerHTML =
                    "Invalid credentials";
                console.log(err);
            });
    });
});
