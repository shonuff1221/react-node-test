const express = require("express");
const router = express.Router();
const studentController = require("./students-controller");
const { handleGetStudentReport } = require("./report-controller");

router.get("", studentController.handleGetAllStudents);
router.post("", studentController.handleAddStudent);
router.get("/:id/report", handleGetStudentReport);
router.get("/:id", studentController.handleGetStudentDetail);
router.post("/:id/status", studentController.handleStudentStatus);
router.put("/:id", studentController.handleUpdateStudent);

module.exports = { studentsRoutes: router };
