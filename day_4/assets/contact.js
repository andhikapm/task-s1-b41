function showData(){
    let showName = document.getElementById('input-name').value
    let showEmail = document.getElementById('input-email').value
    let showNohp = document.getElementById('input-nohp').value
    let showSubject = document.getElementById('input-subject').value
    let showMassage = document.getElementById('input-massage').value
    
    console.log(showName)
    console.log(showEmail)
    console.log(showNohp)
    console.log(showSubject)
    console.log(showMassage)

    if(showName ==''){
        return alert('Need Name')
    }

    if(showEmail ==''){
        return alert('Need Email')
    }
    
    if(showNohp ==''){
        return alert('Need Phone Number')
    }
    
    if(showSubject ==''){
        return alert('Need Subjet?')
    }

    let emailRec = 'andhikapm1031@gmail.com'
    
    let a = document.createElement('a')
    a.href = `mailto:${emailRec}?subject=${showSubject}&body=Hello my name ${showName}, ${showMassage}` 
    a.click()
}