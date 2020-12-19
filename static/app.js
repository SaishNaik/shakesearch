const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const container = document.getElementById("container");
    const rows = [];
    const rowsResult = [];
    for (let result of results) {
      rowsResult.push(`<br><div><pre>${result}</pre><a class="Show More" onclick="myFunction(this)">Read More</a></div>`);
    }
    container.innerHTML = "<div class='result'>Result :" +rowsResult.join(" ");
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);

function myFunction(elem){
  if(elem.innerHTML == "Read More") {
    let spans = elem.parentElement.getElementsByClassName("hide");
    for (let span of spans) {
      span.style.display = "block";
    }
    elem.innerHTML = "Show Less";
  }else{
    let spans = elem.parentElement.getElementsByClassName("hide");
    for (let span of spans) {
      span.style.display = "none";
    }
    elem.innerHTML = "Read More";
  }
}
