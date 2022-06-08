const List = document.getElementById('tasklist')
document.querySelector('#submit').onclick = function(){
    const task =document.createElement("div")
    task.innerText = "hehe"
    List.appendChild(task)
}
