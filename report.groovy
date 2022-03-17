
def functionA()
{
    
    println(jobName)
    job = Jenkins.instance.getJob(jobName)

    //projects = hudson.model.Hudson.instance.getJob('indy-it-JDG-daily').getItems()

    //for( build in projects.getAllJobs())
    //{
    //process build data
    println("All builds: ${job.getBuilds().collect{ it.getNumber()}}");

    //}
}

return this
