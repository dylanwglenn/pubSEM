# Some useful keyboard shortcuts for package authoring:
#
#   Install Package:           'Ctrl + Shift + B'
#   Check Package:             'Ctrl + Shift + E'
#   Test Package:              'Ctrl + Shift + T'


#' Open the pubSEM GUI editor
#'
#' @param fit A lavaan object
#' @param standardized A bool value
#' @param layout_name A string denoting the name of the pubSEM layout to which
#'   a layout should be stored. A layout file stores persistent data for
#'   reproducible diagrams in the current working directoloadry — layout names
#'   should be unique EVEN ACROSS R projects!
#' @returns nothing
#' @export
sem_gui <- function(fit, layout_name, standardized = FALSE) {
    fit_to_json(fit, standardized)

    base_dir <- tools::R_user_dir("pubSEM", which = "data")
    if (Sys.info()['sysname'] == "Windows") {
        gui_exec_path <- system.file("bin", "sem_gui.exe", package = "pubSEM", mustWork = TRUE)
        # run the GUI executable
        system2(gui_exec_path,
                args = c(base_dir, layout_name, "edit"),
                invisible = FALSE #necessary for Windows
        )
    } else {
        gui_exec_path <- system.file("bin", "sem_gui", package = "pubSEM", mustWork = TRUE)
        # run the GUI executable
        system2(gui_exec_path,
                args = c(base_dir, layout_name, "edit")
        )
    }
}

