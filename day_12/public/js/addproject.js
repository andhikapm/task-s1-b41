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
<label for="input-start-date" class="form-label">Start Date</label>
<input type="date" id="input-start-date" class="form-control" name="dataStartDate" value="${formatMin}" min="${formatMin}">`

/*
<label for="input-start-date" class="form-label">Start Date</label>
<input type="date" id="input-start-date" class="form-control" value="${formatMin}" min="${formatMin}">`
*/

/*
<label for="input-start-date">Start Date</label>
<input type="date" id="input-start-date" class="custom-input-start-date" value="${formatMin}" min="${formatMin}">`
*/


//console.log(startDate)

document.getElementById("showEndDate").innerHTML += `
<label for="input-end-date" class="form-label">End Date</label>
<input type="date" id="input-end-date" class="form-control" name="dataEndDate" value="${formatMin}" min="${formatMin}">`

/*
let dataBlog = []

function addBlog(event) {
    event.preventDefault()

    let title = document.getElementById("input-title").value
    let content = document.getElementById("input-content").value
    let image = document.getElementById("input-blog-image").files[0]
    let startDate = document.getElementById("input-start-date").value
    let EndDate = document.getElementById("input-end-date").value
    let datNodeJS = document.getElementById("data-Node-JS").checked
    let datNextJS = document.getElementById("data-Next-JS").checked
    let datReactJS = document.getElementById("data-React-JS").checked
    let datTypeScript = document.getElementById("data-TypeScript").checked

    image = URL.createObjectURL(image)
    console.log(datNodeJS.checked, datNextJS.checked, datReactJS.checked, datTypeScript.checked)

    //calDuration(startDate , EndDate)
    //techList(datNodeJS, datNextJS, datReactJS, datTypeScript)

    let blog = {
        title,
        content,
        image,
        startDate ,
        EndDate,
        datNodeJS,
        datNextJS,
        datReactJS,
        datTypeScript,
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
                    ${calDuration(dataBlog[index].startDate , dataBlog[index].EndDate)}
                </div>
                <div class="parg-content">
                    <p>
                        ${dataBlog[index].content}
                    </p>
                </div>
                <div class="tech-list-group">
                    ${techList(dataBlog[index].datNodeJS, dataBlog[index].datNextJS, dataBlog[index].datReactJS, dataBlog[index].datTypeScript)}
                </div>
                <div class="btn-group">
                    <button class="btn-edit">edit</button>
                    <button class="btn-post">delete</button>
                </div>
            </div>
        </div>
        `
    }

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

    let checkFeb = intStart[0] % 4
    let mon30 = [4, 6, 9, 11]

    //console.log(intStart)
    //console.log(intEnd)
    //console.log(texDuration)

    if(1 >= calYear > 0 && calMonth< 0){

        calMonth = calMonth + 12
        texDuration = calMonth + " Months"
        //console.log(`${calMonth} bulan`)

    } else if(1 >= calMonth > 0 && calDay <0){
        if(checkFeb == 0 && intStart[1] == 2){

            calDay = calDay + 29

        } else if(intStart[1] == 2){

            calDay = calDay + 28

        } else {

            let tanda31 = true

            for(let index = 0; index < 4; index++){

                if(mon30[index] == intStart[1]){
                    calDay = calDay + 30
                    tanda31 = false
                    break 

                }
            }
            if(tanda31 == true){

                calDay = calDay + 31

            }

        }
        texDuration = calDay + " Days"
        //console.log(`${calDay} hari`)

    }else if(calYear > 0){

        texDuration = calYear + " Years"

    }else if(calMonth >0){

        texDuration = calMonth + " Months"

    }else{

        texDuration = calDay + " Days"

    }

    return `Duration : ${texDuration}`
    
}

function techList(dataNode, dataNext, dataReact, dataTypeS){
    let texTechlist = ''
    if(dataNode == true){
        texTechlist = texTechlist + '<i class="fa-brands fa-google"></i>' 
        
    }

    if(dataNext == true){
        texTechlist = texTechlist + '<i class="fa-brands fa-github"></i>' 
        
    }

    if(dataReact == true){
        texTechlist = texTechlist + '<i class="fa-brands fa-windows"></i>' 
        
    }

    if(dataTypeS == true){
        texTechlist = texTechlist + '<i class="fa-brands fa-android"></i>' 
        
    }

    //console.log(texTechlist)

    return texTechlist

}*/