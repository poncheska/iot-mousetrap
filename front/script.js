// let socket = new WebSocket("ws://smart-mousetrap.herokuapp.com/mousetraps/ws");
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
    // if (!response.ok){
    //     let errorMessage = document.createElement('div');
    //     errorMessage.className = "error";
    //     errorMessage.innerHTML = `<strong>${JSON.parse(response).message}</strong>`
    //     document.querySelector("#form-name").after(errorMessage);
    // }
    return response;
}
async function signIn(form) {
    let formdata = new FormData(form);
    if (form.id == "join-form"){
        let joinResponse = await join(formdata);
        let status = joinResponse.ok;
        // let status = join(formdata);
        if (!status){
            if (!document.querySelector("#join-error")){
                // alert(JSON.parse(joinResponse).message);
                let errorMessage = document.createElement('div');
                errorMessage.id = "join-error";
                errorMessage.className = "error";
                // errorMessage.innerHTML = "<strong>fdghjk</strong>";
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
        // let doc =  window.top.document;
        document.querySelector("#open-button").hidden = true;
        let updateButton = document.createElement('button');
        updateButton.id="update-button";
        updateButton.className = "button update-button"
        document.querySelector("#header").after(updateButton);
        updateButton.textContent = "update";
        let table = document.createElement('table');
        table.setAttribute("class", "table");
        table.setAttribute("id", "table");
        table.insertAdjacentHTML("beforeend",`
        <th>name</th>
        <th>status</th>
        <th>last action</th>`);
        let responseMousetraps = await fetch('/mousetraps', {
            headers: {'Content-Type': 'application/json;charset=utf-8',
            'Authorizatioin': `Bearer ${localStorage.getItem('token')}`
        }
        });
        let data = await response.json();
        for (let i = 0; i < data.length; i++){
            table.insertAdjacentHTML("beforeend",`
            <tr>
                <td>${data[i].name}</td>
                <td>${gata[i].status}</td>
                <td>${gata[i].last_trigger}</td>
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
// socket.onmessage = function(event) {
// 	let message = event.data; 
// 	// document.querySelector("#table").
// 	console.log(message);
// 	  }
// document.forms.join.addEventListener("submit", function() {signIn(document.forms.join)});
// document.forms.signIn.addEventListener("submit", function() {signIn(document.forms.signIn)});
function update(){
    let newResponse = await fetch('/mousetraps',{
        headers: {'Content-Type': 'application/json;charset=utf-8',
            'Authorizatioin': `Bearer ${localStorage.getItem('token')}`
        }
    });
    let newData = await newResponse.json();
   
    let existedTable = document.querySelector("#table");
    for (let i = 0; i < existedTable.querySelectorAll('tr').length; i++){
        let DataArray = [newData[i].name, newData[i].status, newData[i].last_trigger];
        for (let j=0; j < existedTable.querySelectorAll('tr')[i].querySelectorAll('td').length; j++){
            existedTable.querySelectorAll('tr')[i].querySelectorAll('td')[j].textContent = DataArray[j];
        }
        
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
document.querySelector("#update-button").addEventListener("click", update);
