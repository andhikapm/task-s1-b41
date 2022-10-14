const realFileBtn = document.getElementById("input-blog-image");
const customBtn = document.getElementById("button-image");
const customTxt = document.getElementById("text-image");

customBtn.addEventListener("click", function() {
  realFileBtn.click();
});

realFileBtn.addEventListener("change", function() {
  if (realFileBtn.value) {
    customTxt.innerHTML = realFileBtn.value.match(
      /[\/\\]([\w\d\s\.\-\(\)]+)$/
    )[1];
  } else {
    customTxt.innerHTML = "";
  }
});

let minTime = new Date()

let minDate = minTime.getDate()

if (minDate <= 9) {
    minDate= "0" + minDate
}
//console.log(minDate)

let minMonth = minTime.getMonth()

if (minMonth <= 9) {
    minMonth = "0" + minMonth
} 

//console.log(minMonth)

let minYear = minTime.getFullYear()
//console.log(minYear)

//19-12-1996 format di butuhkan
let formatMin = `${minYear}-${minMonth}-${minDate}`
//console.log(formatMin)

document.getElementById("showStartDate").innerHTML += `
<label for="input-start-date">Start Date</label>
<input type="date" id="input-start-date" class="custom-input-start-date" value="${formatMin}" min="${formatMin}">`


//console.log(startDate)

document.getElementById("showEndDate").innerHTML += `
<label for="input-end-date">End Date</label>
<input type="date" id="input-end-date" class="custom-input-end-date" value="${formatMin}" min="${formatMin}">`

let dataBlog = []

function addBlog(event) {
    event.preventDefault()

    let title = document.getElementById("input-title").value
    let content = document.getElementById("input-content").value
    let image = document.getElementById("input-blog-image").files[0]
    let startDate = document.getElementById("input-start-date").value
    let EndDate = document.getElementById("input-end-date").value
    /*
    let alpha = startDate.split("-")
    console.log(alpha)
    let beta = parseInt(alpha[0])
    console.log(beta)
    */
    image = URL.createObjectURL(image)
    console.log(image)

    calDuration(startDate , EndDate)

    let blog = {
        title,
        content,
        image,
        startDate ,
        EndDate,
        author: "rangga alfa"
    }
    
    dataBlog.push(blog)
    console.log(dataBlog)

    renderBlog()
}

function renderBlog() {
    document.getElementById("contents").innerHTML = ''

    for (let index = 0; index < dataBlog.length; index++) {
        console.log("test",dataBlog[index])

        document.getElementById("contents").innerHTML += `
        <div class="blog-list-item">
            <div class="blog-image">
                <img src="${dataBlog[index].image}">
            </div>
            <div class="blog-content">
                <h1>
                    <a href="blog-detail.html" target="_blank">
                        ${dataBlog[index].title}
                    </a>
                </h1>
                <div class="detail-blog-content">
                    ${calDuration(dataBlog[index].startDate , dataBlog[index].EndDate)} | ${dataBlog[index].author}
                </div>
                <p>
                    ${dataBlog[index].content}
                </p>
                <div class="btn-group">
                    <button class="btn-edit">Edit Post</button>
                    <button class="btn-post">Post Blog</button>
                </div>
            </div>
        </div>
        `
    }

// ${calDuration(dataBlog[index].startDate , dataBlog[index].EndDate)} | ${dataBlog[index].author}
//${dataBlog[index].startDate} | ${dataBlog[index].author}
}



function calDuration(startTime , endTime){
    let stringStart = startTime.split("-")
    let stringEnd = endTime.split("-")

    //console.log(stringStart)
    //console.log(stringEnd)

    let intStart = []
    let intEnd = []

    for(let index = 0; index < 3; index++){
        intStart[index] = parseInt(stringStart[index])
        intEnd[index] = parseInt(stringEnd[index])
    }

    let calYear = intEnd[0] - intStart[0] // 1
    let calMonth = intEnd[1] - intStart[1] // -9 ---> -9 + 12
    let calDay = intEnd[2] - intStart[2] // -10
    let texDuration = ""

    if(calYear > 0){
        texDuration = calYear + " Years"
    }else if(calMonth >0){
        texDuration = calMonth + " Months"
    }else{
        texDuration = calDay + " Days"
    }

    console.log(intStart)
    console.log(intEnd)
    console.log(texDuration)

    if(calYear > 0 && calMonth< 0){
        calMonth = calMonth + 12
        console.log(calMonth)
    } else if(calMonth > 0 && calDay <0){

    }

    return `Duration : ${texDuration}`
    
}
