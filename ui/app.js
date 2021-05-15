const express = require('express')
const app = express()

var path = require('path');
var bodyParser = require('body-parser');
var invokeev= require("./invokeev");
app.use(express.static(path.join(__dirname, 'public')));
app.use(express.static(path.join(__dirname, 'public', 'select')));
app.set("views", path.join(__dirname, "views"));
//app.set("public", path.join(__dirname, "public"));
app.engine('html', require('ejs').renderFile);
app.set('view engine', 'html');
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
var userID;
var batteryID;
var transactionID;
var stationID;

app.post('/', (req, res) => {
    var name = req.body.name;
    var password = req.body.password;
    if (name == null || name == "") {
        alert("Name can't be blank");
        return false;
    } else if (password.length < 6) {
        alert("Password must be at least 6 characters long.");
        return false;
    } else if (name == "Battery Manufacturer" && password == "12345678") {
        res.sendFile(__dirname+'/views/Battery_Manufacturer.html');
    } else if (name == "Swapping Station" && password == "12345678") {
        res.redirect("select/User_Authentication.html");
    } else if (name == "Electric Vehicle" && password == "12345678") {
        res.redirect("select/Electric_Vehicle.html");
    }
});

app.get('/Add_Battery',(req, res) => {
    res.render('Add_Battery');
});
app.get('/Add_New_Station',(req, res) => {
    res.render('Add_New_Station');
});
app.get('/Add_New_User',(req, res) => {
    res.render('Add_New_User');
});
app.get('/View_Station_Details',(req, res) => {
    res.render('View_Station_Details');
});
app.get('/View_User_Details',(req, res) => {
    res.render('View_User_Details');
});
app.get('/View_Battery_Details',(req, res) => {
    res.render('View_Battery_Details');
});
app.get('/User_New_Battery',(req, res) => {
    res.render('User_New_Battery');
});

app.post('/Add_Battery', async (req, res) => {
    var batId = req.body.batId;
    var batType = req.body.batType;
    var manName = req.body.manName;
    var manAdd = req.body.manAdd;
    var manDate = req.body.manDate;
    var expDate = req.body.expDate;
    var secCert = req.body.secCert;
    var cUser = req.body.cUser;
    await invokeev.createBattery(batId,batType,manName,manAdd,manDate,expDate,secCert,cUser);
    var result = await invokeev.getBattery(batId);
   res.render("battery_details", JSON.parse(result));
});
app.post('/Add_User', async (req, res) => {
    var useId = req.body.useId;    
    console.log("useid:  "+useId);
    var useNm = req.body.useNm;
    var subDets = req.body.subDets;
    var batId = req.body.batId;   
    await invokeev.createevo(useId,useNm,subDets,batId);
    var result = await invokeev.getevo(useId);
     res.render("user_details", JSON.parse(result));
});


app.post('/Add_Station', async (req, res) => {
    var stnId = req.body.stnId;
    var stnNm = req.body.stnNm;
    var stnAdd = req.body.stnAdd;
    var licDets = req.body.licDets;
    await invokeev.createss(stnId,stnNm,stnAdd,licDets);
    var result = await invokeev.getss(stnId);
    res.render("station_details", JSON.parse(result));
});



app.post('/View_Station', async (req, res) => {
    var stnId = req.body.stnId;
    var result = await invokeev.getss(stnId);
    console.log("result: "+result)
    res.render("station_details", JSON.parse(result));
});

app.post('/View_User', async (req, res) => {
    var usrId = req.body.useId;
    var result = await invokeev.getevo(usrId);
    console.log("result: "+result)
    res.render("user_details", JSON.parse(result));
});

app.post('/View_Battery', async (req, res) => {
    var batId = req.body.batId;
    var result = await invokeev.getBattery(batId);
    console.log("result: "+result)
    res.render("battery_details", JSON.parse(result));
});


app.post('/Pay_Bill', async (req, res) => {
    var batteryId = req.body.BatteryID;
    batteryID=batteryId;
    console.log("BatteryID: "+batteryId);
    var batteryUsage = req.body.BatteryUsage;
    transactionID = req.body.TransactionID;
    console.log("TransID: "+transactionID);
    var evOwnerId = userID;
    stationID= req.body.StationID;
    console.log("StationID: "+stationID);
    console.log("Evowner: "+evOwnerId);
    var result = await invokeev.addTransaction(transactionID,batteryId,batteryUsage,evOwnerId);
    console.log("result: "+result);
    invokeev.updateBatteryDetails(batteryID,stationID);
    
    res.render("Pay_Bill", JSON.parse(result));
});

app.post('/Generate_Bill',async (req, res) => {
    var userId = req.body.UserId;
    userID=userId;
    var result = await invokeev.authenticateUser(userId);
    console.log("result: "+result);
    console.log("bool: "+ result.toString()==new String("True"));
    if(result.toString()==new String("True")){
        res.render("Generate_Bill", { userId: userId});
    } else{
        res.send('Invalid User');
    }
    
});

app.get('/confirm', async (req, res) => {
    await invokeev.updatePaidTransaction(transactionID);
    res.render("User_New_Battery", { userID: userID});
});

app.post('/Update_New_Battery', async (req, res) => {
    var batteryId = req.body.BatteryID;
    console.log("UserID: "+ userID);
    await invokeev.updateBatteryCurrentUser(batteryId,userID);
    res.render("Paid");
});

app.listen(3000, (err) => {
    if (err)
        console.log(`Error`);
    console.log(`Running on Port 3000`);
});

