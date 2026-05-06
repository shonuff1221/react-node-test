const asyncHandler = require("express-async-handler");
const { getAllStudents, addNewStudent, getStudentDetail, setStudentStatus, updateStudent } = require("./students-service");

const handleGetAllStudents = asyncHandler(async (req, res) => {
    const students = await getAllStudents(req.query);
    res.json({ students });
});

const handleAddStudent = asyncHandler(async (req, res) => {
    const result = await addNewStudent(req.body);
    res.json(result);
});

const handleUpdateStudent = asyncHandler(async (req, res) => {
    const payload = { ...req.body, id: Number(req.params.id) };
    const result = await updateStudent(payload);
    res.json(result);
});

const handleGetStudentDetail = asyncHandler(async (req, res) => {
    const student = await getStudentDetail(Number(req.params.id));
    res.json(student);
});

const handleStudentStatus = asyncHandler(async (req, res) => {
    const { status } = req.body;
    const result = await setStudentStatus({
        userId: Number(req.params.id),
        reviewerId: req.user.id,
        status
    });
    res.json(result);
});

module.exports = {
    handleGetAllStudents,
    handleGetStudentDetail,
    handleAddStudent,
    handleStudentStatus,
    handleUpdateStudent,
};
