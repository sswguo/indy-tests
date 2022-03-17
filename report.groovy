
def functionA()
{
    projects = Jenkins.instance.getJob('indy-it-JDG-daily').getItems()

    for( build in projects.getAllJobs())
    {
    //process build data
    println(build.getDuration());
    println(build.getTime());

    }
}
