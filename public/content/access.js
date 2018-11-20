window.addEventListener("load", () => {
    document.querySelector("button#search").addEventListener("click", () => {
        const username = document.querySelector("#username").value;
        const accessLvl = document.querySelector("#accessLvl").value;
        const employeeId = document.querySelector("#employeeId").value;

        const baseUrl = `${window.location.protocol}//${window.location.host}`;
        fetch(`${baseUrl}/api/content/searchAccess`, {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({
                username,
                accessLvl: accessLvl === "" ? 0 : parseInt(accessLvl),
                employeeId: employeeId === "" ? 0 : parseInt(employeeId)
            })
        })
            .then(res => {
                if (res.ok) {
                    return res.json();
                } else {
                    throw res.statusText;
                }
            })
            .then(json => {
                tableCreate(json);
            })
            .catch(err => {
                document.querySelector("#result").innerHTML = err;
                console.log(err);
            });
    });
});

const tableCreate = json => {
    const createCell = (value, td) => {
        td.appendChild(document.createTextNode(value));
    };

    const tbl = document.createElement("table");
    tbl.id = "result";
    tbl.className = "table";
    const thead = document.createElement('thead');
    const tr = document.createElement('tr');
    let th;

    const columns = ["Employee ID", "Username", "Password", "Access Level", "Computer Access", "IP", "MAC"];

    for (let column of columns) {
        th = document.createElement('th');
        th.appendChild(document.createTextNode(column));
        tr.appendChild(th);
    }

    thead.appendChild(tr);
    tbl.appendChild(thead);
    const body = document.createElement('tbody');

    for (let i = 0; i < json.length; i++) {
        const tr = body.insertRow();
        createCell(json[i]["employeeId"], tr.insertCell());
        createCell(json[i]["username"], tr.insertCell());
        createCell(json[i]["password"], tr.insertCell());
        createCell(json[i]["accessLvl"], tr.insertCell());
        createCell(json[i]["computerAccess"], tr.insertCell());
        createCell(json[i]["IP"], tr.insertCell());
        createCell(json[i]["MAC"], tr.insertCell());
    }

    tbl.appendChild(body);

    // build a reference to the existing node to be replaced
    const sp2 = document.getElementById("result");
    const parentDiv = sp2.parentNode;

    // replace existing node sp2 with the new span element sp1
    parentDiv.replaceChild(tbl, sp2);
};
