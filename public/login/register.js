window.addEventListener("load", () => {
    document.querySelector("button#login").addEventListener("click", () => {
        const username = document.querySelector("#username").value;
        const password = document.querySelector("#password").value;
        const baseUrl = `${window.location.protocol}//${window.location.host}`;
        fetch(`${baseUrl}/api/register`, {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({ username, password })
        })
            .then(() => {
                document.querySelector("#result").innerHTML = "Registered";
            })
            .catch(err => {
                document.querySelector("#result").innerHTML =
                    "Invalid credentials";
                console.log(err);
            });
    });
});
