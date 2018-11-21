window.addEventListener("load", () => {
    document.querySelector("button#login").addEventListener("click", () => {
        const username = document.querySelector("#username").value;
        const password = document.querySelector("#password").value;
        const baseUrl = `${window.location.protocol}//${window.location.host}`;
        const mode = document.querySelector("#mode").value;
        fetch(`${baseUrl}/api/login-${mode}`, {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({ username, password })
        })
            .then(res => {
                if (res.ok) {
                    window.location.assign(`${baseUrl}/content/${mode}.html`);
                } else {
                    document.querySelector("#result").innerHTML =
                        "Invalid credentials";
                }
            })
            .catch(err => {
                document.querySelector("#result").innerHTML =
                    "Invalid credentials";
                console.log(err);
            });
    });
});
