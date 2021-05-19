let modalSignIn = document.querySelector("#modal-sign-in"),
modalJoin = document.querySelector("#modal-join"),
modalOverlay = document.querySelector("#modal-overlay"),
closeButtonSignIn = document.querySelector("#close-button-sign-in"),
closeButtonJoin = document.querySelector("#close-button-join"),
openButtonSignIn = document.querySelector("#open-button-sign-in");
openButtonJoin = document.querySelector("#open-button-join");
async function join(formdata) {
    let response = await fetch("/org/sign-up", {
        method: 'POST',
        headers: {'Content-Type': 'application/json;charset=utf-8'
    },
    body: JSON.stringify({"name":formdata.get("Email"), "pass":formdata.get("Password")})
    });
    return response;
}
async function signIn(form) {
    let formdata = new FormData(form);
    if (form.id == "join-form"){
        let joinResponse = await join(formdata);
        let status = joinResponse.ok;
        if (!status){
            if (!document.querySelector("#join-error")){
                let errorMessage = document.createElement('div');
                errorMessage.id = "join-error";
                errorMessage.className = "error";
                let json = await joinResponse.json();
                errorMessage.innerHTML = `<strong>${json.message}</strong>`;
                document.querySelector("#form-name-join").after(errorMessage);
            }
            return false;
         };
    }
    let response = await fetch("/org/sign-in", {
        method: 'POST',
        headers: {'Content-Type': 'application/json;charset=utf-8'
    },
    body: JSON.stringify({"name":formdata.get("Email"), "pass":formdata.get("Password")})
    });
    let jsonResponse = await response.json();
    if (response.ok){
        localStorage.setItem('token', jsonResponse.token);
        if (form.id == "join-form"){
            modalJoin.classList.toggle("closed");}
        else{
            modalSignIn.classList.toggle("closed");
        }
        modalOverlay.classList.toggle("closed"); 
        document.querySelector("#open-button").hidden = true;
        let updateButton = document.createElement('button');
        updateButton.id="update-button";
        updateButton.className = "button update-button";
        updateButton.textContent = "update";
        document.querySelector("#user-buttons").append(updateButton);
        let logOutButton = document.createElement('button');
        logOutButton.id="log-out-button";
        logOutButton.className = "button log-out-button";
        logOutButton.textContent = "log out";
        document.querySelector("#update-button").after(logOutButton);
        // let jsonUpdateButton = JSON.stringify(updateButton,['id','className','textContent']);
        // localStorage.setItem('updateButton', jsonUpdateButton);
        let table = document.createElement('table');
        table.setAttribute("class", "table");
        table.setAttribute("id", "table");
        table.insertAdjacentHTML("beforeend",`
        <th>name</th>
        <th>status</th>
        <th>last action</th>`);
        let responseMousetraps = await fetch('/mousetraps', {
            headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
        });
        let data = await responseMousetraps.json();
        localStorage.setItem('data', JSON.stringify(data));
        for (let i = 0; i < data.length; i++){
            table.insertAdjacentHTML("beforeend",`
            <tr>
                <td>${data[i].name}</td>
                <td>${data[i].status}</td>
                <td>${data[i].last_trigger}</td>
            </tr>`)
        }
        document.querySelector("#header").after(table);
    }
    else{
        if (!document.querySelector("#sign-in-error")){
            let errorMessage = document.createElement('div');
            errorMessage.className = "error";
            errorMessage.id = "sign-in-error";
            errorMessage.innerHTML = `<strong>${jsonResponse.message}</strong>`
            document.querySelector("#form-name-sign-in").after(errorMessage);}
    }
    return false
}
function buttonDefinition(modal, modalOverlay, closeButton, openButton) {
	closeButton.addEventListener("click", function() {
		modal.classList.toggle("closed");
		modalOverlay.classList.toggle("closed");
        //modal.querySelector("#Email").value = "";
        //modal.querySelector("#Password").value = "";
	});

    openButton.addEventListener("click", function() {
        modal.classList.toggle("closed");
        modalOverlay.classList.toggle("closed");
    });
    }
    buttonDefinition(modalSignIn, modalOverlay, closeButtonSignIn, openButtonSignIn);
    buttonDefinition(modalJoin, modalOverlay, closeButtonJoin, openButtonJoin);
async function update(){
    let newResponse = await fetch('/mousetraps',{
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
    });
    let newData = await newResponse.json();
    localStorage.setItem('data', JSON.stringify(newData));
    let existedTable = document.querySelector("#table");
    for (let i = 0; i < existedTable.querySelectorAll('tr').length - 1; i++){
        let dataArray = [newData[i].name, newData[i].status, newData[i].last_trigger];
        for (let j=0; j < existedTable.querySelectorAll('tr')[i+1].querySelectorAll('td').length; j++){
            existedTable.querySelectorAll('tr')[i+1].querySelectorAll('td')[j].textContent = dataArray[j];
        }
        
    }
}
function savingChanges(){
    if (localStorage.length){
        document.querySelector("#open-button").hidden = true;
        let updateButton = document.createElement('button');
        updateButton.id="update-button";
        updateButton.className = "button update-button";
        updateButton.textContent = "update";
        document.querySelector("#user-buttons").append(updateButton);
        let table = document.createElement('table');
        table.setAttribute("class", "table");
        table.setAttribute("id", "table");
        table.insertAdjacentHTML("beforeend",`
        <th>name</th>
        <th>status</th>
        <th>last action</th>`);
        let data = JSON.parse(localStorage.getItem('data'));
        for (let i = 0; i < data.length; i++){
            table.insertAdjacentHTML("beforeend",`
            <tr>
                <td>${data[i].name}</td>
                <td>${data[i].status}</td>
                <td>${data[i].last_trigger}</td>
            </tr>`)
        }
        document.querySelector("#header").after(table);
        let logOutButton = document.createElement('button');
        logOutButton.id="log-out-button";
        logOutButton.className = "button log-out-button";
        logOutButton.textContent = "log out";
        document.querySelector("#update-button").after(logOutButton);
    }
}
document.forms.join.addEventListener("submit", function(event) {
	event.preventDefault();
	signIn(document.forms.join);
    event.currentTarget.submit();
});
document.forms.signIn.addEventListener("submit", function(event) {
    event.preventDefault();
	signIn(document.forms.signIn);
    event.currentTarget.submit();
});
savingChanges()
document.querySelector("#update-button").addEventListener("click", update);
document.querySelector("#update-button").addEventListener("mousedown", function(event) {this.style.background = '#5a1200'});
document.querySelector("#update-button").addEventListener("mouseup", function(event) {this.style.background = '#ff3300'});
document.querySelector(".button").addEventListener("mouseover", function(event) {this.style.background = '#a52100'});
document.querySelector(".button").addEventListener("mouseout", function(event) {this.style.background = '#ff3300'});
document.querySelector("#log-out-button").addEventListener("click", function(event) {this.style.background = '#5a1200'; localStorage.clear(); window.location.reload()});
