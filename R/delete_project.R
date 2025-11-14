#' Delete an existing project
#'
#' @param project_name A string denoting the name of the pubSEM project to
#'   delete.
#' @returns nothing
#' @export
delete_project <- function(project_name) {
    base_dir <- tools::R_user_dir("pubSEM", which = "data")
    project_file <- file.path(base_dir, paste0(project_name, ".json"))

    file.remove(project_file)
    print(paste("Successfully deleted the following project:", project_name))
}
